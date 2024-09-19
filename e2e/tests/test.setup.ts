import {expect} from '@playwright/test';
import {Client} from '@e2e-support/server';
import {test as setup} from '@e2e-support/test_fixture';
import {starterTemplatePluginId} from '@support/constant';

setup('ensure server has license', async ({pw}) => {
    const {adminClient} = await pw.getAdminClient();
    expect(await ensureLicense(adminClient)).toBe(true);
});

setup('ensure plugin is enabled', async ({pw}) => {
    const {adminClient} = await pw.getAdminClient();

    const pluginStatus = await adminClient.getPluginStatuses();
    const plugins = await adminClient.getPlugins();

    for (const pluginId of [starterTemplatePluginId]) {
        const isInstalled = pluginStatus.some(({plugin_id}) => plugin_id === pluginId);

        if (!isInstalled) {
            console.log(`${pluginId} is not installed. Related visual test will fail.`);
            continue;
        }

        const isActive = plugins.active.some(({id}) => id === pluginId);

        if (!isActive) {
            await adminClient.enablePlugin(pluginId);
            console.log(`${pluginId} is installed and has been activated.`);
        } else {
            console.log(`${pluginId} is installed and active.`);
        }
    }
});

async function ensureLicense(adminClient: Client) {
    try {
        const currentLicense = await adminClient.getClientLicenseOld();

        if (currentLicense?.IsLicensed === 'true') {
            return true;
        }

        await requestTrialLicense(adminClient);

        const trialLicense = await adminClient.getClientLicenseOld();
        return trialLicense?.IsLicensed === 'true';
    } catch (error) {
        console.error('Error ensuring license', error);
        return false;
    }
}

async function requestTrialLicense(adminClient: Client) {
    try {
        // @ts-expect-error This may fail requesting for trial license
        await adminClient.requestTrialLicense({
            receive_emails_accepted: true,
            terms_accepted: true,
            users: 100,
        });
    } catch (e) {
        console.error('Failed to request trial license', e);
        throw e;
    }
}
