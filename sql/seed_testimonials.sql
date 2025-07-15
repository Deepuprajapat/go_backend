-- SQL script to seed testimonials data into static_site_data
-- Replace YOUR_STATIC_SITE_DATA_ID with the actual id value

UPDATE static_site_data
SET testimonials = '[
  {
    "created_date": 1728545088910,
    "description": "Now the search is over finally, I would like to thanks Mr. Aditya Sharma for all the help and support. He offered Ace Divino to us and in just one visit we decided to book our dream home. His marketing skills and attention to detail made the process easy and lucrative for us. He help us for further processing too which was seemless and efficient. That’s not the easy lift but with your efforts all settled on time. Thanks once again, will love to have more happy experiences with you in future deals.",
    "is_approved": true,
    "name": "Priyanka Khatri",
    "rating": 5,
    "type": "Client ",
    "updated_date": 1737622816291,
    "user_id": 1
  },
  {
    "created_date": 1728545097942,
    "description": "I am satisfied with Invest Mango team. They supervised and made the whole process super easy. The team is very friendly and professional at the same time.",
    "is_approved": true,
    "name": "Chaudhary Gaurav Malik",
    "rating": 5,
    "type": "Client ",
    "updated_date": 1737622373063,
    "user_id": 1
  },
  {
    "created_date": 1728545221043,
    "description": " Invest Mango is excellent in his service and deal. We are in touch with Mr.Varun last 2 years, Finally after long time, I have done deal with help of Varun. Thanks for your nice effort. ",
    "is_approved": true,
    "name": "Saurabh Kumar Agrawal",
    "rating": 5,
    "type": "Client ",
    "updated_date": 1737622869537,
    "user_id": 1
  },
  {
    "created_date": 1737622649174,
    "description": "I would like to thanks Invest Mango for helping us out and giving the best supervision during the whole process of purchasing this property Thank you, this is what I was looking for",
    "is_approved": true,
    "name": "George Gomes",
    "rating": 5,
    "type": "Client",
    "updated_date": null,
    "user_id": 1
  }
]'
-- WHERE id = 'YOUR_STATIC_SITE_DATA_ID';
