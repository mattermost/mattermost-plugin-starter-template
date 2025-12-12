import {test, expect} from '@e2e-support/test_fixture';

test('Can access plugin settings', async ({pw, pages}) => {
    // # Log in as admin
    const {adminUser} = await pw.initSetup();
    const {page} = await pw.testBrowser.login(adminUser);

    // # Visit system console
    const systemConsolePage = new pages.SystemConsolePage(page);
    await systemConsolePage.goto();
    await systemConsolePage.toBeVisible();

    // # Go to PLUGINS > Plugin Starter Template
    await systemConsolePage.page.getByRole('link', {name: 'Plugin Starter Template'}).click();

    // # Enable the plugin
    await page.getByTestId('PluginSettings.PluginStates.com+mattermost+plugin-starter-template.Enabletrue').check();
    await page.getByTestId('saveSetting').click();

    // * Verify that the plugin is active and ready to use
    await expect(page.getByTestId('saveSetting')).toBeDisabled();
});
