import streamlit as st
import json
import time
import datetime
from dummy import get_dummy_response
from utils import render_offer, render_all_offers

ENDPOINT_URL = "http://localhost:80/api/offers"

st.set_page_config(layout="wide")

# INIT SESSION STATES

if "current_state" not in st.session_state:
    st.session_state["current_state"] = 0 ## 0 = BEGIN, 1 = SEARCH, 2

if "region_num_dict" not in st.session_state:
    with open("region_num_dict.json", "r", encoding="utf-8") as file:
        print("Read regions")
        d = json.load(file)
    st.session_state["region_num_dict"] = d

if "num_region_dict" not in st.session_state:
    with open("num_region_dict.json", "r", encoding="utf-8") as file:
        print("Read regions")
        d = json.load(file)
    st.session_state["num_region_dict"] = d



if "current_page_index" not in st.session_state:
    st.session_state["current_page_index"] = 0

if "page_size" not in st.session_state:
    st.session_state["page_size"] = 20    

if "offers" not in st.session_state:
    st.session_state["offers"] = []

min_seat_num = 1
only_vollkasko = 0
min_price = 0
max_price = 0
car_type = None
min_free_km = 0
refresh_button = False
prev = False
next = False

st.header("Go Api Go - The Best Rental Car Service", divider="red")

if st.session_state["current_state"] == 0:
    st.write(" ")
    st.write(" ")
    st.write(" ")
    st.write(" ")
    st.write(" ")
    st.write(" ")
    st.write(" ")
    st.write(" ")
    st.write(" ")
c1, c2,c3,c4, c5, c6, c7, c8 = st.columns([1,0.6,0.5,0.6,0.5,0.4,0.8,0.5], vertical_alignment="bottom")
with c1:
    region_options = st.session_state["region_num_dict"].keys()
    region_options = sorted(region_options)
    selected = st.selectbox("Region", 
                region_options,
                index=None,
                placeholder="Choose a Region")
    if selected is None:
        region_id = -1
    else:
        region_id = st.session_state["region_num_dict"][selected]
with c2:
    begin_date = st.date_input("Time Range Begin", value="today", label_visibility='visible')
with c3:
    begin_time = st.time_input("Start Time Input",value="now", label_visibility="hidden")

with c4:
    end_date = st.date_input("Time Range End", value="today", label_visibility="visible")
with c5:
    end_time = st.time_input("End Time Input",value="now", label_visibility="hidden")

with c6:
    amount_days = st.number_input(label="Days",value=10, min_value=1, step=1, placeholder="Choose the amount of days (24h)")

with c7:
    order = st.selectbox(label="Order", options=["Price Ascending", "Price Descending"], index=0)
    if order == "Price Ascending":
        order = "price-asc"
    elif order == "Price Descending":
        order = "price-desc"
    else:
        order = None
with c8:
    st.write("")
    st.write("")
    search_button = st.button("Search", use_container_width=True, type="primary")

# Show Big Menu
if st.session_state["current_state"] == 1:
    c1, c2 = st.columns([0.001,6])
    with st.sidebar:
        st.write(" ")

        st.write("**Filter Options**")
        min_seat_num = st.number_input("Min. Seat Number", min_value=1,step=1, max_value=50, value=1, placeholder="Min. Free Seats")
        only_vollkasko = st.toggle("Only Vollkasko", value=False)

        min_price, max_price = st.slider(label="Price Range", min_value=0,
                    max_value=1000,
                    value=(0, 1000))
        
        car_type = st.selectbox("Car Type", options=["small", "sports", "luxury", "family"],index=None, placeholder="Choose type of car")

        min_free_km = st.number_input("Min. free Kilometers", min_value=0, step=50, value=0)

        refresh_button = st.button("Refresh", type="primary", use_container_width=True)
    
    # ALL THE OFFERS
    with c2:
        
        render_all_offers(st.session_state["offers"], st.session_state["page_size"])
# Choose Page 
    st.write(" ")
    _, c1, c2, c3,c4,c5,_ = st.columns([1.1,1, 0.75,0.75,0.8,0.75,1.7], vertical_alignment="center")
    with c1:
        st.write("Ergebnisse pro Seite: ")
    with c2:
        page_size = st.selectbox("Num. Results per Page", options=[20,50,100], index=0, label_visibility="collapsed")
        if not page_size == st.session_state["page_size"]:
            st.session_state["page_size"] = page_size
    with c3:
        prev = st.button("Previous",use_container_width=True, disabled=(st.session_state["current_page_index"] == 0))

    with c4:
        current_page = st.session_state["current_page_index"]
        st.write(f"**Current Page: {current_page}**")
    with c5:
        current_page = st.session_state["current_page_index"] + 1
        next = st.button("Next", use_container_width=True)



## BACKEND ##################
def wait_for_response():
    c1, c2, c3 = st.columns([3,1,3])
    with c2:
        with st.spinner():
            response = get_dummy_response()
            return response


def perform_search():
    page_index = st.session_state["current_page_index"]
    page_size = st.session_state["page_size"]
    print(f"""Perform search for Region {region_id}, 
          Begin Date {begin_date}, 
          Begin Time {begin_time}, 
          End Date {end_date}, 
          End Time {end_time}, 
          Amount Days {amount_days}, 
          Order {order},
          Page {page_index}, 
          PageSize {page_size},
          MinSeatNum {min_seat_num},
          OnlyVollkasko {only_vollkasko},
          Price Range ({min_price},{max_price}),
          Car Types {car_type},
          MinFreeKM {min_free_km}""")
    
    combined_start_time = datetime.datetime.combine(begin_date, begin_time)
    combined_end_time = datetime.datetime.combine(end_date, end_time)
    unixtime_start = time.mktime(combined_start_time.timetuple())
    unixtime_end = time.mktime(combined_end_time.timetuple())

    
    send_data = {"regionID":region_id,
                "timeRangeStart":unixtime_start,
                "timeRangeEnd":unixtime_end,
                "numberDays": amount_days,
                "sortOrder":order, 
                "page":page_index,
                "pageSize":page_size,
                "priceRangeWidth":100,      # NOT IMPLEMENTED
                "minFreeKilometerWidth":100 # NOT IMPLEMENTED
                }
    if min_seat_num > 1:
        send_data.update({"minNumberSeats":min_seat_num})
    if min_price > 0:
        send_data.update({"minPrice":min_price})
    if max_price > 0:
        send_data.update({"maxPrice":max_price})
    if not car_type is None:
        send_data.update({"carType":car_type}) 
    if only_vollkasko:
        send_data.update({"onlyVollkasko":only_vollkasko})
    if min_free_km > 0:
        send_data.update({"minFreeKilometer":min_free_km})

    print("Send Data", send_data)

    # response = requests.post(ENDPOINT_URL, data={})
    response = wait_for_response()
    st.session_state["offers"] = response
    
if search_button:

    if st.session_state["current_state"] == 0:
        st.session_state["current_state"] = 1
        perform_search()
        st.rerun()
    else:
        perform_search()
        st.rerun()

if refresh_button:
    perform_search()
    st.rerun()

if prev:
    st.session_state["current_page_index"] -= 1
    perform_search()
    st.rerun()

if next:
    st.session_state["current_page_index"] += 1
    perform_search()
    st.rerun()







