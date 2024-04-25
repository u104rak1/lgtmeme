# Endpoints

### Health check endpoint
| Endpoint                   | Method                      | Path                                 |
| :------------------------- | :-------------------------- | :----------------------------------- |
| Health Check               | `HEAD`                      | /api/health                          |

### Auth server api endpoints
| Endpoint                   | Method                      | Path                                 |
| :------------------------- | :-------------------------- | :----------------------------------- |
| Authorization              | `GET`                       | /auth-api/authorize                  |
| Access Token               | `POST`                      | /auth-api/token                      |
| Public key                 | `GET`                       | /auth-api/jwks                       |
| Login                      | `POST`                      | /auth-api/login                      |
| Logout                     | `GET`                       | /auth-api/logout                     |

### Auth server view endpoints
| Endpoint                   | Method                      | Path                                 |
| :------------------------- | :-------------------------- | :----------------------------------- |
| Login                      | `GET`                       | /login                               |

### Client server endpoints
| Endpoint                   | Method                      | Path                                 |
| :------------------------- | :-------------------------- | :----------------------------------- |
| Callback                   | `GET`                       | /client-api/admin/callback           |
| Create image bff           | `POST`                      | /client-api/images                   |
| Get images bff             | `GET`                       | /client-api/images                   |
| Update image bff           | `PATCH`                     | /client-api/images/:image_id         |
| Delete image bff           | `Delete`                    | /client-api/images/:image_id         |


### Client server view endpoints
| Endpoint                   | Method                      | Path                                 |
| :------------------------- | :-------------------------- | :----------------------------------- |
| Home                       | `GET`                       | /                                    |
| New                        | `GET`                       | /new                                 |
| Privacy policy             | `GET`                       | /privacy-policy                      |
| Terms of service           | `GET`                       | /terms-of-service                    |
| Admin                      | `GET`                       | /admin                               |

### Resource server endpoints
| Endpoint                   | Method                      | Path                                 |
| :------------------------- | :-------------------------- | :----------------------------------- |
| Create image               | `POST`                      | /resource-api/images                 |
| Get images                 | `GET`                       | /resource-api/images                 |
| Update image               | `PATCH`                     | /resource-api/images/:image_id       |
| Delete image               | `Delete`                    | /resource-api/images/:image_id       |
