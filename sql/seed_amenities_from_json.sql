-- SQL script to seed amenities data into static_site_data from JSON export
-- This script groups amenities by category and limits to max 2 amenities per category

UPDATE static_site_data
SET categories_with_amenities = '{
  "categories": {
    "CONVENIENCE": [
      {
        "icon": "https://image.gopropify.com/amenityImages/cafe_coffee_bar.svg",
        "value": "Cafe Coffee Bar"
      },
      {
        "icon": "https://image.gopropify.com/amenityImages/high_speed_escalators.svg",
        "value": "High Speed Escalators"
      }
    ],
    "SPORTS": [
      {
        "icon": "https://image.gopropify.com/amenityImages/basketball.svg",
        "value": "Basket Ball"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/gymnasium.svg",
        "value": "gymnasium"
      }
    ],
    "GARDEN": [
      {
        "icon": "https://image.gopropify.com/amenityImages/large_green_area.svg",
        "value": "Large Green Area"
      }
    ],
    "ENTERTAINMENT": [
      {
        "icon": "https://image.gopropify.com/amenityImages/multiplex.svg",
        "value": "Multiplex"
      }
    ],
    "FACILITIES": [
      {
        "icon": "https://image.gopropify.com/amenityImages/main_entrance.png",
        "value": "Main Entrance"
      },
      {
        "icon": "https://image.gopropify.com/amenityImages/24x7_water_supply.svg",
        "value": "24x7 Water Supply"
      }
    ],
    "SAFETY": [
      {
        "icon": "https://image.investmango.com/images/update/img/icon/24x7_security.svg",
        "value": "24x7_security"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/cctv_video_surveillance.svg",
        "value": "cctv_video_surveillance"
      }
    ],
    "LEISURE": [
      {
        "icon": "https://image.investmango.com/images/update/img/icon/sauna.svg",
        "value": "sauna"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/party_hall.png",
        "value": "party_hall"
      }
    ],
    "ENVIRONMENT": [
      {
        "icon": "https://image.investmango.com/images/update/img/icon/rain_water_harvesting.svg",
        "value": "rain_water_harvesting"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/sewage_treatment_plant.svg",
        "value": "sewage_treatment_plant"
      }
    ],
    "LIFESTYLE": [
      {
        "icon": "https://image.investmango.com/images/update/img/icon/Pergola.png",
        "value": "pergola"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/gazebo.png",
        "value": "gazebo"
      }
    ]
  }
}'
WHERE is_active = true; 