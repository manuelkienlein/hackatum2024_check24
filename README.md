# HackaTUM 2024 - Check24 Challenge for high-performance Car Sharing Plattform

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

This repository contains the code for our solution to the HackaTUM Check24 challenge. The challenge required us to build a high-performance REST API for managing and querying car rental offers, along with a web interface for easy interaction. We’ve also included a log parser written in Python to help us debug the API and monitor its performance.

## Description

This project aims to solve the problem of managing and querying car rental offers with high performance, while also providing a user-friendly interface for interacting with the system. The core components of the solution are:

- **REST API**: Built with Go (using the Fiber framework) to handle car rental offers.
- **Database**: PostgreSQL is used to store the rental offer data.
- **Frontend**: A beautiful web UI built using Python and Streamlit.
- **Log Parser**: A Python-based log parser for debugging and analyzing API calls.

### Features
- **Add Car Rental Offers**: Add rental offers with details like price, car type, region, availability, and more.
- **Query Offers**: Filter and sort offers based on criteria like price, car type, region, and availability.
- **High-Performance REST API**: Built with Go and optimized for speed, handling high traffic and large datasets.
- **Web Interface**: Built using Python with Streamlit, providing an easy way to interact with the API.
- **Log Parsing**: A Python tool to parse API logs, making it easier to debug and monitor the system’s behavior.

### Technologies
Backend:
- Go with Fiber framework for the REST API.
- PostgreSQL as the database for storing car rental offers.

Frontend:
- Python with Streamlit for the web interface.
