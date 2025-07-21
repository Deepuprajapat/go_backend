-- SQL script to seed amenities data into static_site_data
-- Replace YOUR_STATIC_SITE_DATA_ID with the actual id value

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
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/atm.svg",
        "value": "ATM"
      }
    ],
    "SPORTS": [
      {
        "icon": "https://image.gopropify.com/amenityImages/basketball.svg",
        "value": "Basket Ball"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/gymnasium.svg",
        "value": "Gymnasium"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/swimming_pool.svg",
        "value": "Swimming Pool"
      }
    ],
    "SAFETY": [
      {
        "icon": "https://image.investmango.com/images/update/img/icon/24x7_security.svg",
        "value": "24x7 Security"
      },
      {
        "icon": "https://image.investmango.com/images/update/img/icon/cctv_video_surveillance.svg",
        "value": "CCTV Video Surveillance"
      }
    ]
  }
}'
-- WHERE id = 'YOUR_STATIC_SITE_DATA_ID'; 