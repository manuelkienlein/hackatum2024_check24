import streamlit as st
import json
import collections

st.set_page_config(layout="wide")

# INIT SESSION STATES

if "current_state" not in st.session_state:
    st.session_state["current_state"] = 0 ## 0 = BEGIN, 1 = SEARCH, 2

if "region_num_dict" not in st.session_state:
    with open("region_num_dict.json", "r", encoding="utf-8") as file:
        print("Read regions")
        d = json.load(file)
    st.session_state["region_num_dict"] = d

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
c1, c2,c3,c4, c5, c6, c7, c8 = st.columns([1,0.6,0.5,0.6,0.5,0.4,0.8,0.5])
with c1:
    region_options = st.session_state["region_num_dict"].keys()
    print(region_options)
    region_options = sorted(region_options)
    print(region_options)
    selected = st.selectbox("Region", 
                region_options,
                index=None,
                placeholder="Choose a Region")
    if selected is None:
        region_id = -1
    else:
        region_id = st.session_state["region_num_dict"][selected]
with c2:
    st.date_input("Time Range Begin", value="today", label_visibility='visible')
with c3:
    st.time_input("Start Time Input",value="now", label_visibility="hidden")

with c4:
    st.date_input("Time Range End", value="today", label_visibility="visible")
with c5:
    st.time_input("End Time Input",value="now", label_visibility="hidden")

with c6:
    st.number_input(label="Days",value=10, min_value=1, step=1, placeholder="Choose the amount of days (24h)")

with c7:
    st.selectbox(label="Order", options=["Price Ascending", "Price Descending"], index=0)
with c8:
    st.write("")
    st.write("")
    search_button = st.button("Search", use_container_width=True, type="primary")

# Show Big Menu
if st.session_state["current_state"] == 1:
    c1, c2 = st.columns([1,6])
    with c1:
        st.write(" ")

        st.write("**Filter Options**")
        st.number_input("Min. Seat Number", min_value=1,step=1, max_value=10, value=1, placeholder="Min. Free Seats")
        st.toggle("Only Vollkasko", value=False)

        st.slider(label="Price Range", min_value=0,
                    max_value=500,
                    value=(100, 400))
        
        st.multiselect("Car Type", options=["small", "sports", "luxury", "family"], placeholder="Choose type of car")

        st.number_input("Min. free Kilometers", min_value=0, step=50, value=0)
    
    # Choose Page part
    st.write(" ")
    _, c2, c3 = st.columns([2,1,4])
    with c2:
        st.selectbox("Num. Results per Page", options=[20,50,100], index=0)
    with c3:
        st.segmented_control("Pages",["Next Page"])
    

if search_button and st.session_state["current_state"] == 0:
    st.session_state["current_state"] = 1
    st.rerun()