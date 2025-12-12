### Setup your environment

In order to get your environment set up to run [Playwright](https://playwright.dev) tests, you can run `./setup-environment`, or run equivalent commands for your current setup.

What this script does:

-   Navigate to the folder above `mattermost-plugin-starter-template`
-   Clone `mattermost` (if it is already cloned there, please have a clean git index to avoid issues with conflicts)
-   Install webapp dependencies - `cd mattermost/webapp && npm i`
-   Install Playwright test dependencies - `cd ../e2e-tests/playwright && npm i`
-   Install plugin e2e dependencies - `cd ../../../mattermost-plugin-starter-template/e2e && npm i`
-   Build and deploy plugin - `make deploy`

---

### Run the tests

Start Mattermost server:

-   `cd <path>/mattermost/server`
-   `make test-data`
-   `make run-server`

Run test:

-   `cd <path>/mattermost-plugin-starter-template/e2e`
-   `npm run test` to run in multiple projects such as `chrome`, `firefox` and `ipad`.
-   `npm run test -- --project=chrome` to run in specific project such as `chrome`.

To see the test report:

-   `cd <path>/mattermost-plugin-starter-template/e2e`
-   `npm run show-report` and navigate to link provided

To see test screenshots:

-   `cd <path>/mattermost-plugin-starter-template/e2e/screenshots`

### Have questions or wanted to connect?

If you have any questions or need assistance, feel free to join and start a discussion at the [QA: Contributors channel](https://community.mattermost.com/core/channels/qa-contributors). It's a good place to connect with Mattermost staff and the wider community.
