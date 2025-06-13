import {PlaywrightTestConfig} from '@playwright/test';
import dotenv from 'dotenv';
import testConfig from '@e2e-test.playwright-config';
dotenv.config({path: `${__dirname}/.env`});

// Configuration override for plugin tests
testConfig.testDir = __dirname + '/tests';
testConfig.outputDir = __dirname + '/test-results';

const projects = testConfig.projects?.map((p) => ({...p, dependencies: ['setup']})) || [];
testConfig.projects = [{name: 'setup', testMatch: /test.setup.ts/} as PlaywrightTestConfig].concat(projects);

export default testConfig;
