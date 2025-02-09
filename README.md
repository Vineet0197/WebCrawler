## **Web Crawler for E-commerce Product URLs**
A high-performance, scalable web crawler written in GoLang that extracts product URLs from multiple e-commerce websites using XPath, Kafka, BFS, and ML-based dynamic pattern detection.

## **Features**
- Multi-threaded crawling using Goroutines
- Kafka-based queue for managing pending URLs
- XPath-based product URL extraction
- ML model for dynamic product page detection
- robots.txt compliance
- Structured JSON output

## **System Architecture**
Below is the high-level architecture of the **Web Crawler**, showing how components interact.

![System Architecture](https://kroki.io/plantuml/svg/eNqFlE9v2kAQxe_7KUY-VERqYvVfDlZVxaWQtJBi2aRUqnpYzIBX2F5rd0yIqn73zi4mxahSuSAzv30z-97gG0vSUFuVghSVCAtcwtDIxxINxCYvFGFOrUEhZE7awINFIxo-onLVyJogiJPPkKHZoQlAWuDHfn0i11vpS5N4PIn7xWOrhTZbNBYGtxpS3ZKq0V74Q8M0XkxHaf-Y0UtN9or2BN9kqVaSR_N0Ovs4m_fZu_k8gWGpsCaPuOc-8T2RVEAije3ukMRpdt7yfsoi0lq1Vh11P-0TX7LZV8h4ErlBD2TzWRrfjoRwpsHlB-dNBMksm0OYu4sL5x3_7o3hSrsslS1gpSupanhIp1b4kmM6HyIY6tq2FZ5SMOj8uxBiY3TbwBgpL-DF0R0UwJ9Owql5n1irwHwLf930mK_1WsZlqR_DT1g_CaxXxx4Z36FR9QYSo3O09ryHMzrqJrmbs1uu7tM41U6R16v2AOS6Jo5JnCsd8ohgtCfDawiN0auWv71Djj0A_9L1yLne_TQ6hvn0LNZwbAc1nuRUaaxK4vx2zsp-a2_GiW4XeASZ3GEPBU7K7YcQtSYEozYFgV4_r8ispaYlWGtTSYrEL54jwL2smhKvcl0FEfzwowUFUWOjMDwphl2j8NXrN8HL_2Nv310HTP0Uv90FwA3Urej7y8OKpmgb3jKEgf9_uowz95rAFS_YDR9y74s_cStRyQ==)

## **Tech Stack**
- Golang (Echo Framework)
- Kafka (Streaming Domains parallely)
- Docker and Docker Compose (Containerized environment)

### **API Endpoints**
#### Submit Domains for Crawling
- Method: `POST /crawl`
- Request Body:
```json
[
  "https://example.com",
  "https://amazon.com"
]
```
- Description: Adds domains to Kafka queue for processing.
#### Extracted Product URL Output
Extracted URLs are stored in a JSON file (product_data.json):
```json
{
  "example.com": [
    "https://example.com/product/123",
    "https://example.com/product/456"
  ],
  "amazon.com": [
    "https://amazon.com/item/abc",
    "https://amazon.com/item/xyz"
  ]
}
```
## **Installation and Setup**
### Prerequisites
1. Ensure you have the following installed:
- Docker & Docker Compose
- Go 1.21+
2. Clone the Repository
```sh
git clone https://github.com/yourusername/web-crawler.git
cd web-crawler
```
3. Build and Run with Docker
```sh
docker-compose up -d --build
```
4. Stop the Services
```sh
docker-compose down
```

### Author
- Developed By: Vineet Aggarwal
- GitHub: Vineet0197