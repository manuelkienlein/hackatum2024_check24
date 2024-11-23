import streamlit as st
from datetime import datetime

LUX_CAR_IMAGE_PATH = "resources/luxury.png"
LUX_CAR_IMAGE_PATH_V2 = "resources/v2/luxury_icon.jpg"
SPORT_CAR_IMAGE_PATH = "resources/sportscar.png"
SPORT_CAR_IMAGE_PATH_V2 = "resources/v2/sportscar_icon.jpg"
FAMILY_CAR_IMAGE_PATH = "resources/familycar.png"
FAMILY_CAR_IMAGE_PATH_V2 = "resources/v2/familycar_icon.jpg"
SMALL_CAR_IMAGE_PATH = "resources/smallcar.png"
SMALL_CAR_IMAGE_PATH_V2 = "resources/v2/smallcar_icon.jpg"


car_type_to_image = { 
    "small":SMALL_CAR_IMAGE_PATH_V2,
    "sports":SPORT_CAR_IMAGE_PATH_V2,
    "luxury":LUX_CAR_IMAGE_PATH_V2,
    "family":FAMILY_CAR_IMAGE_PATH_V2
}

def render_all_offers(offers, limit_index):

    st.write(" ")
    for index, offer in enumerate(offers[0:limit_index]):
        render_offer(index+1, offer)



def render_offer(index, offer):
    num_region_dict = st.session_state["num_region_dict"]
    region = num_region_dict[str(offer["RegionID"])]
    number_seats = str(offer["NumberSeats"])
    price = int(offer["Price"])
    car_type = str(offer["CarType"])
    car_image = car_type_to_image[car_type]
    free_km = str(offer["FreeKilometers"])
    only_vollkasko = offer["HasVollkasko"]
    start_timestamp = str(offer["StartTimestamp"])
    end_timestamp = str(offer["EndTimestamp"])
    start_date = datetime.strptime(start_timestamp, "%Y-%m-%dT%H:%M:%SZ")
    end_date = datetime.strptime(end_timestamp, "%Y-%m-%dT%H:%M:%SZ")
    formatted_date = f"{start_date.strftime('%d/%m/%Y')}-{end_date.strftime('%d/%m/%Y')}"
    if only_vollkasko:
        insurance_symbol = "✅"
    else:
        insurance_symbol = "❌"

    padding, col1,padding2, col2, col3, col4 = st.columns([0.001,1,0.3,2,2,2])
    with col1:
        st.write(" ")
        st.image(car_image, width=120)
    with col2:
        st.write("  ")
        st.write("  ")
        st.write(f"{formatted_date}")
        st.write(f"Price: **{float(price/100):.2f}€**")
    with col3:
        st.write("  ")
        st.write("  ")
        st.write(f"Region: {region}")
        st.write(f"Number Seats: {number_seats}")
    with col4:
        st.write("  ")
        st.write("  ")
        st.write(f"Vollkasko: {insurance_symbol}")
        st.write(f"Free Range: {free_km} km")
