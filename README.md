<br />
<div align="center">
  <h3 align="center">Short-URL</h3>

  <p align="center">
    A URL shortener system
    <br />
    <br />
    <a href="#demo">View Demo</a>
    ·
    <a href="#system-design-overview">System Design Overview</a>
  </p>
</div>

<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#system-design-overview">System Design Overview</a></li>
  </ol>
</details>

## About The Project

This project build based on karanpratapsingh [url-shortener](https://github.com/karanpratapsingh/system-design#url-shortener) system design guideline.

#### What is a URL Shortener?

A URL shortener service creates an alias or a short URL for a long URL. Users are redirected to the original URL when they visit these short links.

For example, the following long URL can be changed to a shorter URL.

**Long URL**: [https://karanpratapsingh.com/courses/system-design/url-shortener](https://karanpratapsingh.com/courses/system-design/url-shortener)

**Short URL**: [https://bit.ly/3I71d3o](https://bit.ly/3I71d3o)

## Getting Started

This is an example of how you may give instructions on setting up your project locally. To get a local copy up and running follow these simple example steps.

### Prerequisites

  * [https://nodejs.org/en/download/](https://nodejs.org/en/download/)
  * [https://www.docker.com/](https://www.docker.com/)

### Installation

- Clone the repo
   ```sh
   git clone https://github.com/aaydin-tr/short-url.git
   ```
- Setup environment variables
    - Back-End
    ```sh
    cd server/
    cp .env.dev .env
    ``` 
    Set up necessary environment variables such as MONGO_URL, MONGO_USERNAME, REDIS_URL, etc.

    - Front-End
    ```sh
    cd web/
    cp .env.dev .env
    ``` 
    If you did not change the server default port you do not need to change anything. Otherwise, you need to set correct `VITE_API_URL` URL

- Use `docker-compose` to start back-end server
   ```sh
   cd server/ && docker-compose up -d
   ```
- Npm run dev
  After all containers start running, you can start web server.
  ```sh
  npm run dev
  ``` 

With default configs, back-end server will run on [http://localhost:8090](http://localhost:8090), and front-end server will run on [http://localhost:5173/](http://localhost:5173/)

## System Design Overview

<p align="center">
  <img src="https://user-images.githubusercontent.com/50546655/192112444-70525473-df7e-4ec4-a28c-282d1334cf5d.png">
</p>

### Sequence Diagram
  ```mermaid
  sequenceDiagram
      actor Client
      Client->>Server: GET /shorturl
      activate Server
      Server->>Cache: get shorturl
      activate Cache
      alt Get original url form cache
          Cache-->>Server: original-url.com
          Server-->>Client: 302 Found / Location: original-url.com
      else Get original url form database
          Cache-->>Server: not found
          Server->>Database: get shorturl
          activate Database
          alt Short Url exist in database
              Database-->>Server: original-url.com
              Server-)Cache: set original-url.com
              deactivate Cache
              Server-->>Client: 302 Found / Location: original-url.com
          else  Short Url not exist in database
              Database-->>Server: not found
              Server-->>Client: 404 Not Found
          end
          deactivate Database
      end
      deactivate Server
  ```
