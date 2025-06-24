# Properties API Documentation

## Overview
This documentation covers the Properties API endpoints for managing real estate properties in the system. The API provides comprehensive property data including property details, web cards, images, pricing, and more.

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

### 1. GET Properties
Retrieve a list of properties with detailed information.

**Endpoint:** `GET /properties`

#### Headers
```http
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

#### Response
```json
{
    "status": 200,
    "data": {
        "properties": [
            {
                "id": "string",
                "name": "string",
                "property_images": ["string"],
                "web_cards": {
                    "property_details": {
                        "built_up_area": {
                            "value": "string"
                        },
                        "sizes": {
                            "value": "string"
                        },
                        "floor_number": {
                            "value": "string"
                        },
                        "configuration": {
                            "value": "string"
                        },
                        "possession_status": {
                            "value": "string"
                        },
                        "balconies": {
                            "value": "string"
                        },
                        "covered_parking": {
                            "value": "string"
                        },
                        "bedrooms": {
                            "value": "string"
                        },
                        "property_type": {
                            "value": "string"
                        },
                        "age_of_property": {
                            "value": "string"
                        },
                        "furnishing_type": {
                            "value": "string"
                        },
                        "rera_number": {
                            "value": "string"
                        },
                        "facing": {
                            "value": "string"
                        },
                        "bathrooms": {
                            "value": "string"
                        }
                    },
                    "why_choose_us": {
                        "image_urls": ["string"],
                        "usp_list": ["string"]
                    },
                    "know_about": {
                        "description": "string"
                    },
                    "video_presentation": {
                        "title": "string",
                        "video_url": "string"
                    },
                    "location_map": {
                        "description": "string",
                        "google_map_link": "string"
                    }
                },
                "pricing_info": {
                    "price": "string"
                },
                "property_rera_info": {
                    "rera_number": "string"
                },
                "meta_info": {
                    "title": "string",
                    "description": "string",
                    "keywords": "string",
                    "canonical": "string"
                },
                "developer_id": "string",
                "location_id": "string",
                "project_id": "string",
                "is_featured": boolean,
                "is_deleted": boolean,
                "created_at": "string",
                "updated_at": "string"
            }
        ]
    }
}
```

#### Response Fields Description

##### Core Property Information
- `id`: Unique identifier for the property
- `name`: Name of the property
- `property_images`: Array of property image URLs

##### Web Cards
###### Property Details
- `built_up_area`: Built-up area of the property
- `sizes`: Property sizes
- `floor_number`: Floor number
- `configuration`: Property configuration (e.g., 2BHK, 3BHK)
- `possession_status`: Current possession status
- `balconies`: Number of balconies
- `covered_parking`: Covered parking details
- `bedrooms`: Number of bedrooms
- `property_type`: Type of property
- `age_of_property`: Age of the property
- `furnishing_type`: Furnishing status
- `rera_number`: RERA registration number
- `facing`: Property facing direction
- `bathrooms`: Number of bathrooms

###### Why Choose Us
- `image_urls`: Array of images showcasing property highlights
- `usp_list`: List of unique selling points

###### Know About
- `description`: Detailed description of the property

###### Video Presentation
- `title`: Title of the video
- `video_url`: URL of the property video

###### Location Map
- `description`: Location description
- `google_map_link`: Google Maps link

##### Pricing Information
- `pricing_info.price`: Property price

##### RERA Information
- `property_rera_info.rera_number`: RERA registration number

##### Meta Information
- `title`: SEO title
- `description`: SEO description
- `keywords`: SEO keywords
- `canonical`: Canonical URL

##### Relationships
- `developer_id`: ID of the developer
- `location_id`: ID of the location
- `project_id`: ID of the associated project

##### Status Flags
- `is_featured`: Whether the property is featured
- `is_deleted`: Whether the property is deleted

##### Timestamps
- `created_at`: Creation timestamp
- `updated_at`: Last update timestamp

#### Status Codes
- `200 OK`: Successfully retrieved properties
- `401 Unauthorized`: Invalid or missing authentication token
- `403 Forbidden`: Insufficient permissions to access properties
- `500 Internal Server Error`: Server error while processing the request

### 2. PATCH Property Details
Update only the property details of an existing property.

**Endpoint:** `PATCH /properties/{property_id}`

#### Parameters
- `property_id` (path, required): Unique identifier of the property to update

#### Headers
```http
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

#### Request Body
```json
{
    "web_cards": {
        "property_details": {
            "built_up_area": {
                "value": "1200 sq ft"
            },
            "sizes": {
                "value": "1200 sq ft"
            },
            "floor_number": {
                "value": "3"
            }
        }
    }
}
```

#### Response
```json
{
    "status": 200,
    "data": {
        "property": {
            // Returns the complete updated property object
            // Only property_details will be modified
            // Other fields will remain unchanged
        }
    },
    "message": "Property details updated successfully"
}
```

#### Status Codes
- `200 OK`: Property details updated successfully
- `400 Bad Request`: Invalid request body or validation errors
- `401 Unauthorized`: Invalid or missing authentication token
- `403 Forbidden`: Insufficient permissions to update the property
- `404 Not Found`: Property not found
- `500 Internal Server Error`: Server error while processing the request

#### Notes
1. This endpoint allows updating only the property details section
2. All fields within property_details are optional - only include the ones you want to update
3. Existing values for fields not included in the request will remain unchanged
4. The response includes the complete property object with updated property details
5. The `updated_at` timestamp will be automatically updated by the system 

### 3. DELETE Property
Delete a property (supports both soft delete and hard delete).

**Endpoint:** `DELETE /properties/{property_id}`

#### Parameters
- `property_id` (path, required): Unique identifier of the property to delete
- `hard_delete` (query, optional): Boolean flag for permanent deletion (default: false)

#### Headers
```http
Authorization: Bearer <your-jwt-token>
```

#### Examples
```bash
# Soft delete (default)
DELETE /properties/123

# Hard delete (permanent)
DELETE /properties/123?hard_delete=true
```

#### Response (Soft Delete)
```json
{
    "status": 200,
    "data": {
        "property": {
            "id": "string",
            "name": "string",
            "is_deleted": true,
            "deleted_at": "2024-03-21T10:00:00Z",
            "updated_at": "2024-03-21T10:00:00Z"
            // Other property fields will be included but omitted for brevity
        }
    },
    "message": "Property soft deleted successfully"
}
```

#### Response (Hard Delete)
```json
{
    "status": 204,
    "data": null,
    "message": "Property permanently deleted successfully"
}
```

#### Status Codes
- `200 OK`: Property soft deleted successfully
- `204 No Content`: Property permanently deleted successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Invalid or missing authentication token
- `403 Forbidden`: Insufficient permissions to delete the property
- `404 Not Found`: Property not found
- `409 Conflict`: Property cannot be deleted (e.g., has active dependencies)
- `500 Internal Server Error`: Server error while processing the request

#### Notes
1. Soft delete (default behavior):
   - Sets `is_deleted` flag to true
   - Updates `deleted_at` timestamp
   - Property remains in database but won't appear in regular GET requests
   - Can be restored later if needed

2. Hard delete (when `hard_delete=true`):
   - Permanently removes the property from the database
   - Cannot be undone
   - All associated data will be deleted
   - Use with caution

3. Dependencies:
   - If the property has active dependencies (e.g., active bookings, linked documents), the delete operation may fail
   - Check the response message for specific details about why a delete operation failed

4. Authorization:
   - Requires appropriate permissions to delete properties
   - Some properties may have additional deletion restrictions based on their status 