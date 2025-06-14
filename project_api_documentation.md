# Project API Documentation

## Overview
This documentation covers the Project Detail API endpoints for managing real estate projects in the system. The API supports full CRUD operations with comprehensive project data including metadata, web cards, images, floor plans, amenities, and more.

## Base URL
```
https://api.investmango.com/v1
```

## Authentication
All endpoints require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Content Type
All requests and responses use JSON format:
```
Content-Type: application/json
```

---

## API Endpoints

### 1. GET Project Detail
Retrieve detailed information about a specific project.

**Endpoint:** `GET /projects/{project_id}`

#### Parameters
- `project_id` (path, required): Unique identifier of the project

#### Response
```json
{
    "data": {
        "project_id": 1,
        "project_name": "Project Name",
        "search_context": ["search_context_1", "search_context_2"],
        "meta_info": {
            "title": "Project Title",
            "description": "Project Description",
            "keywords": "Project Keywords",
            "canonical": "https://example.com/project",
            "project_schema": "<script type=\"application/ld+json\">\n{\n  \"@context\": \"https://schema.org/\",\n  \"@type\": \"Product\",\n  \"name\": \"ACE Divino\",\n  \"image\": \"https://image.investmango.com/images/img/ace-divino/ace-divino-greater-noida-west.webp\",\n  \"description\": \"ACE Divino Sector 1, Noida Extension: Explore prices, floor plans, payment options, location, photos, videos, and more. Download the project brochure now!\",\n  \"brand\": {\n    \"@type\": \"Brand\",\n    \"name\": \"Ace Group of India\"\n  },\n  \"offers\": {\n    \"@type\": \"AggregateOffer\",\n    \"url\": \"https://www.investmango.com/ace-divino\",\n    \"priceCurrency\": \"INR\",\n    \"lowPrice\": \"18800000\",\n    \"highPrice\": \"22500000\"\n  }\n}\n</script>"
        },
        "web_cards": {
            "images_url": [
                "https://example.com/image.jpg",
                "https://example.com/image.jpg"
            ],
            "project_info": {
                "name": "Project Name",
                "description": "Project Description",
                "area": "Project Area",
                "logo_url": "https://example.com/logo.jpg",
                "min_price": "Min Price",
                "max_price": "Max Price"
            },
            "rera_list": [
                {
                    "rera_number": "RERA Number",
                    "phase": "RERA Phase", 
                    "status": "RERA Status",
                    "rera_qr": "RERA QR"
                }
            ],
            "project_details": {
                "area": {
                    "value": "Project Area"
                },
                "sizes": {
                    "value": "Sizes"
                },
                "units": {
                    "value": "Project Units"
                },
                "launch_date": {
                    "value": "Launch Date"
                },
                "possession_date": {
                    "value": "Possession Date"
                },
                "total_towers": {
                    "value": "Total Towers"
                },
                "total_floors": {
                    "value": "Total Floors"
                },
                "status": {
                    "value": "Project Status"
                },
                "type": {
                    "value": "Property Type"
                }
            },
            "why_to_choose": {
                "image_urls": [
                    "https://example.com/image.jpg"
                ],
                "usp_list": [
                    {
                        "icon": "https://example.com/icon.jpg",
                        "html_content": "Why to Choose"
                    }
                ],
            },
            "know_about": {
                "description": "HTML content",
                "download_link": "https://example.com/download"
            },
            "floor_plan": {
                "description": "Floor Plan",
                "products": [
                    {
                        "title": "2 BHK Apartment + 2 Toilets",
                        "flat_type": "2BHK",
                        "price": "Price",
                        "building_area": "Building Area",
                        "image": "https://example.com/image.jpg",
                    }
                ]
            },
            "price_list": {
                "description": "Price List",
                "bhk_options_with_prices": [
                    {
                        "bhk_option": "2BHK",
                        "size": "1000 sq.ft.",
                        "price": "1000000"
                    }
                ]
            },
            "payment_plans": {
                "description": "Payment Plans",
                "plans": [
                    {
                        "name": "Plan 1",
                        "details": "Plan Details"
                    }
                ]
            },
            "amenities": {
                "title": "Amenities",
                "amenities_with_categories": {
                    "category_1": [
                        {
                        "icon": "https://example.com/icon.jpg",
                        "value": "Amenity 1"
                        }
                    ],
                    "category_2": [
                        {
                        "icon": "https://example.com/icon.jpg",
                        "value": "Amenity 2"
                        }
                    ]
                }
            },
            "video_presentation": {
                "description": "Video Presentation",
                "url": "https://example.com/video"
            },
            "location_info": {
                "google_map_link": "https://example.com/map",
                "city": "Noida",
                "longitude": "123.456",
                "latitude": "78.901"
            },
            "site_plan": {
                "description": "HTML content",
                "image_url": "https://example.com/image.jpg"
            },
            "about": {
                "description": "HTML content",
                "developer_logo_url": "https://example.com/logo.jpg",
                "establishment_year": "2021",
                "total_properties": "12",
                "contact_details": {
                    "name": "Contact Name",
                    "address": "Contact Address",
                    "phone": "Contact Phone",
                    "booking_link": "https://example.com/booking"
                }
            },
            "faqs": [
                {
                    "question": "Question 1",
                    "answer": "Answer 1"
                },
                {
                    "question": "Question 2",
                    "answer": "Answer 2"
                }
            ]
        },
        "developer_info": {
            "name": "Developer Name",
            "logo_url": "https://example.com/logo.jpg",
            "established_year": "2021",
            "total_projects": "12"
        },
        "location_info": {
            "google_map_link": "https://example.com/map",
            "city": "Noida",
            "longitude": "123.456",
            "latitude": "78.901"
        },
        "is_premium": true,
        "is_priority": true,
        "is_featured": true,
        "deleted_at": null,
        "created_at": "2021-01-01T00:00:00Z",
        "updated_at": "2021-01-01T00:00:00Z"
    },
    "status": 200
}
```

#### Status Codes
- `200 OK`: Project found and returned successfully
- `404 Not Found`: Project does not exist
- `401 Unauthorized`: Invalid or missing authentication token

---

### 2. POST Create Project
Create a new project with complete details.

**Endpoint:** `POST /projects`

#### Request Body
```json
{
 "project_name": "Project Name",
 "project_url": "https://example.com/project",
 "project_type": "Project Type",
 "locality": "Project Locality",
 "project_city": "Project City",    
 "developer_id": 1,
}
```

#### Response
```json
{
    "project_id": 1,
}
```

#### Status Codes
- `201 Created`: Project created successfully
- `400 Bad Request`: Validation errors in request data
- `401 Unauthorized`: Invalid or missing authentication token
- `403 Forbidden`: Insufficient permissions to create projects

---

### 3. PATCH Update Project
Update specific fields of an existing project (partial update).

**Endpoint:** `PATCH /projects/{project_id}`

#### Parameters
- `project_id` (path, required): Unique identifier of the project to update

#### Request Body (All fields optional)
```json
{
    "project_name": "Updated Project Name",
    "search_context": ["updated_context_1", "updated_context_2"],
    "meta_info": {
        "title": "Updated Project Title",
        "description": "Updated Project Description"
    },
    "web_cards": {
        "project_info": {
            "min_price": "Updated Min Price",
            "max_price": "Updated Max Price"
        },
        "amenities": {
            "title": "Updated Amenities"
        }
    },
    "is_premium": false,
    "is_featured": true
}
```

#### Response
```json
{
    "status": 200,
    "data": {
        "project_id": 1,
        "project_name": "Updated Project Name",
        "search_context": ["updated_context_1", "updated_context_2"],
        "meta_info": {
            "title": "Updated Project Title",
            "description": "Updated Project Description",
            "keywords": "Project Keywords",
            "project_url": "https://example.com/project",
            "project_schema": "<script type=\"application/ld+json\">\n{\n  \"@context\": \"https://schema.org/\",\n  \"@type\": \"Product\",\n  \"name\": \"ACE Divino\",\n  \"image\": \"https://image.investmango.com/images/img/ace-divino/ace-divino-greater-noida-west.webp\",\n  \"description\": \"ACE Divino Sector 1, Noida Extension: Explore prices, floor plans, payment options, location, photos, videos, and more. Download the project brochure now!\",\n  \"brand\": {\n    \"@type\": \"Brand\",\n    \"name\": \"Ace Group of India\"\n  },\n  \"offers\": {\n    \"@type\": \"AggregateOffer\",\n    \"url\": \"https://www.investmango.com/ace-divino\",\n    \"priceCurrency\": \"INR\",\n    \"lowPrice\": \"18800000\",\n    \"highPrice\": \"22500000\"\n  }\n}\n</script>"
        },
        "web_cards": {
            "images": [
                {
                    "order": 1,
                    "url": "https://example.com/image.jpg"
                },
                {
                    "order": 2,
                    "url": "https://example.com/image.jpg"
                }
            ],
            "project_info": {
                "project_name": "Updated Project Name",
                "project_description": "Project Description",
                "project_area": "Project Area",
                "project_group": "Project Group",
                "min_price": "Updated Min Price",
                "max_price": "Updated Max Price"
            },
            "rera_info": {
                "rera_number": "RERA Number",
                "rera_phase": "RERA Phase",
                "rera_status": "RERA Status",
                "rera_qr": "RERA QR"
            },
            "project_details": {
                "project_area": {
                    "value": "Project Area"
                },
                "sizes": {
                    "value": "Sizes"
                },
                "project_units": {
                    "value": "Project Units"
                },
                "launch_date": {
                    "value": "Launch Date"
                },
                "possession_date": {
                    "value": "Possession Date"
                },
                "total_towers": {
                    "value": "Total Towers"
                },
                "total_floors": {
                    "value": "Total Floors"
                },
                "project_status": {
                    "value": "Project Status"
                },
                "property_type": {
                    "value": "Property Type"
                }
            },
            "why_to_choose": {
                "images": [
                    {
                        "order": 1,
                        "url": "https://example.com/image.jpg"
                    }
                ],
                "usps": [
                    {
                        "icon": "https://example.com/icon.jpg",
                        "value": "Why to Choose"
                    }
                ],
                "expert_link": "https://example.com/expert",
                "booking_link": "https://example.com/booking"
            },
            "know_about": {
                "html_text": "HTML Text",
                "download_link": "https://example.com/download"
            },
            "floor_plan": {
                "description": "Floor Plan",
                "products": [
                    {
                        "title": "2 BHK Apartment + 2 Toilets",
                        "configuration": "2BHK",
                        "price": "Price",
                        "area": "1000 sq.ft.",
                        "image": "https://example.com/image.jpg",
                        "expert_link": "https://example.com/expert",
                        "brochure_link": "https://example.com/brochure"
                    }
                ]
            },
            "price_list": {
                "description": "Price List",
                "bhk_options_with_prices": [
                    {
                        "configuration": "2BHK",
                        "size": "1000 sq.ft.",
                        "price": "1000000"
                    }
                ]
            },
            "payment_plans": {
                "title": "Payment Plans",
                "plans": [
                    {
                        "plan_name": "Plan 1",
                        "plan_details": "Plan Details"
                    }
                ]
            },
            "amenities": {
                "title": "Updated Amenities",
                "amenities_with_categories": {
                    "category_1": {
                        "icon": "https://example.com/icon.jpg",
                        "value": "Amenity 1"
                    },
                    "category_2": {
                        "icon": "https://example.com/icon.jpg",
                        "value": "Amenity 2"
                    }
                }
            },
            "video_presentation": {
                "description": "Video Presentation",
                "url": "https://example.com/video"
            },
            "location_info": {
                "description": "HTML content",
                "google_map_link": "https://example.com/map"
            },
            "site_plan": {
                "html_content": "HTML content",
                "image": "https://example.com/image.jpg"
            },
            "about": {
                "description": "About",
                "logo_url": "https://example.com/logo.jpg",
                "establishment_year": "2021",
                "total_properties": "12",
                "html_content": "HTML content",
                "contact_details": {
                    "name": "Contact Name",
                    "address": "Contact Address",
                    "phone": "Contact Phone",
                    "booking_link": "https://example.com/booking"
                }
            },
            "faqs": [
                {
                    "question": "Question 1",
                    "answer": "Answer 1"
                },
                {
                    "question": "Question 2",
                    "answer": "Answer 2"
                }
            ]
        },
        "is_premium": false,
        "is_priority": true,
        "is_featured": true,
        "deleted_at": null,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T10:30:00Z"
    }
}
```

#### Status Codes
- `200 OK`: Project updated successfully
- `400 Bad Request`: Validation errors in request data
- `404 Not Found`: Project does not exist
- `401 Unauthorized`: Invalid or missing authentication token
- `403 Forbidden`: Insufficient permissions to update this project

---

### 4. DELETE Project
Delete a project (supports both soft delete and hard delete).

**Endpoint:** `DELETE /projects/{project_id}`

#### Parameters
- `project_id` (path, required): Unique identifier of the project to delete
- `hard_delete` (query, optional): Boolean flag for permanent deletion (default: false)

#### Examples
```bash
# Soft delete (default)
DELETE /projects/1

# Hard delete (permanent)
DELETE /projects/1?hard_delete=true
```

#### Response (Soft Delete)
```json
{
    "status": 200,
    "data": {
        "project_id": 1,
        "project_name": "Project Name",
        "search_context": ["search_context_1", "search_context_2"],
        "meta_info": {
            "title": "Project Title",
            "description": "Project Description",
            "keywords": "Project Keywords",
            "project_url": "https://example.com/project"
        },
        "web_cards": {
            "images": [
                {
                    "order": 1,
                    "url": "https://example.com/image.jpg"
                }
            ],
            "project_info": {
                "project_name": "Project Name",
                "project_description": "Project Description",
                "project_area": "Project Area",
                "project_group": "Project Group",
                "min_price": "Min Price",
                "max_price": "Max Price"
            },
            "rera_info": {
                "rera_number": "RERA Number",
                "rera_phase": "RERA Phase",
                "rera_status": "RERA Status",
                "rera_qr": "RERA QR"
            }
        },
        "is_premium": true,
        "is_priority": true,
        "is_featured": true,
        "deleted_at": "2024-01-01T12:00:00Z",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T12:00:00Z"
    },
    "message": "Project soft deleted successfully"
}
```

#### Response (Hard Delete)
```json
{
    "status": 204,
    "data": null,
    "message": "Project permanently deleted successfully"
}
```

#### Status Codes
- `200 OK`: Project soft deleted successfully
- `204 No Content`: Project permanently deleted successfully
- `404 Not Found`: Project does not exist
- `409 Conflict`: Project already deleted or has dependencies
- `401 Unauthorized`: Invalid or missing authentication token
- `403 Forbidden`: Insufficient permissions to delete this project

---

## Data Models

### Project Structure
The project data model consists of several nested objects:

#### Core Fields
- `project_id`: Unique identifier (auto-generated)
- `project_name`: Name of the project
- `search_context`: Array of search keywords
- `is_premium`: Boolean flag for premium status
- `is_priority`: Boolean flag for priority status
- `is_featured`: Boolean flag for featured status
- `deleted_at`: Timestamp for soft deletion (null if active)
- `created_at`: Creation timestamp
- `updated_at`: Last modification timestamp

#### Meta Information
- `title`: SEO title
- `description`: SEO description
- `keywords`: SEO keywords
- `project_url`: Project landing page URL
- `project_schema`: JSON-LD structured data

#### Web Cards
Complex nested structure containing:
- **Images**: Project gallery with ordering
- **Project Info**: Basic project details
- **RERA Info**: Real Estate Regulatory Authority details
- **Project Details**: Technical specifications
- **Why to Choose**: USPs and highlights
- **Know About**: Detailed descriptions
- **Floor Plan**: Layout plans and pricing
- **Price List**: BHK options with prices
- **Payment Plans**: Available payment schemes
- **Amenities**: Categorized amenities list
- **Video Presentation**: Promotional videos
- **Location Info**: Map and location details
- **Site Plan**: Master plan information
- **About**: Developer information
- **FAQs**: Frequently asked questions

---

## Error Handling

### Error Response Format
```json
{
    "status": 400,
    "message": "Error description",
    "error_code": "ERROR_CODE",
    "errors": [
        {
            "field": "field_name",
            "message": "Field-specific error message"
        }
    ]
}
```

### Common Error Codes
- `PROJECT_NOT_FOUND`: Project with specified ID does not exist
- `PROJECT_ALREADY_DELETED`: Attempting to delete an already deleted project
- `CONSTRAINT_VIOLATION`: Cannot delete due to existing dependencies
- `VALIDATION_ERROR`: Request data validation failed
- `UNAUTHORIZED`: Authentication required or token invalid
- `FORBIDDEN`: Insufficient permissions for the operation

---

## Usage Examples

### cURL Examples

#### Get Project
```bash
curl -X GET \
  'https://api.investmango.com/v1/projects/1' \
  -H 'Authorization: Bearer your-jwt-token' \
  -H 'Content-Type: application/json'
```

#### Create Project
```bash
curl -X POST \
  'https://api.investmango.com/v1/projects' \
  -H 'Authorization: Bearer your-jwt-token' \
  -H 'Content-Type: application/json' \
  -d '{
    "project_name": "New Project",
    "meta_info": {
      "title": "New Project Title",
      "description": "New Project Description"
    },
    "is_premium": true
  }'
```

#### Update Project
```bash
curl -X PATCH \
  'https://api.investmango.com/v1/projects/1' \
  -H 'Authorization: Bearer your-jwt-token' \
  -H 'Content-Type: application/json' \
  -d '{
    "project_name": "Updated Project Name",
    "is_featured": true
  }'
```

#### Delete Project (Soft)
```bash
curl -X DELETE \
  'https://api.investmango.com/v1/projects/1' \
  -H 'Authorization: Bearer your-jwt-token'
```

#### Delete Project (Hard)
```bash
curl -X DELETE \
  'https://api.investmango.com/v1/projects/1?hard_delete=true' \
  -H 'Authorization: Bearer your-jwt-token'
```

---

## Best Practices

### Performance Considerations
1. **Pagination**: Use pagination for listing endpoints to manage large datasets
2. **Field Selection**: Consider implementing field selection to reduce payload size
3. **Caching**: Implement appropriate caching strategies for frequently accessed data
4. **Rate Limiting**: Implement rate limiting to prevent abuse

### Security Guidelines
1. **Authentication**: Always validate JWT tokens
2. **Authorization**: Implement role-based access control
3. **Input Validation**: Validate all input data thoroughly
4. **Sanitization**: Sanitize HTML content in web cards
5. **HTTPS**: Always use HTTPS in production

### Data Integrity
1. **Validation**: Implement comprehensive validation rules
2. **Transactions**: Use database transactions for multi-step operations
3. **Soft Delete**: Prefer soft delete for data recovery capabilities
4. **Audit Trail**: Maintain audit logs for all modifications

---

## Rate Limits
- **GET requests**: 1000 requests per hour per API key
- **POST/PATCH requests**: 100 requests per hour per API key
- **DELETE requests**: 50 requests per hour per API key

## Support
For API support and questions, contact: api-support@investmango.com

---

*Last updated: 2024-01-01* 