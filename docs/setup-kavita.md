# Setup Kavita

### Steps

1. Make sure the following environment variables are set in your `.env` file:
    ```properties
    MEDIA_DIR=/path/to/your/media
    ```

2. Get Kavita running in your Docker environment by following these steps:
    ```sh
    docker compose up -d kavita
    ```

3. Access Kavita by navigating to `http://localhost:5000` in your web browser.

4. Follow the Kavita setup wizard to configure your library and start managing your media collection.

5. Once setup is complete, you can access your media library and enjoy your content!

6. For le-grimoire integration, get the Kavita API key from the [settings page](http://10.65.127.83:5000/settings#clients) by following this [guide](https://wiki.kavitareader.com/guides/api/) and add it to your `.env` file:
    ```properties
    KAVITA_API_KEY=your_kavita_api_key
    ```

#### Shutdown Kavita

```sh
docker compose down kavita
```
