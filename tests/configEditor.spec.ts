import { test, expect } from '@grafana/plugin-e2e';
import { UpdownDataSourceOptions, UpdownSecureJsonData } from '../src/types';

test('smoke: should render config editor', async ({ createDataSourceConfigPage, readProvisionedDataSource, page }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await createDataSourceConfigPage({ type: ds.type });
  await expect(page.getByText('API URL')).toBeVisible();
  await expect(page.getByText('API Key')).toBeVisible();
});

test('should show API URL field with default value', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<UpdownDataSourceOptions, UpdownSecureJsonData>({ fileName: 'datasources.yml' });
  await createDataSourceConfigPage({ type: ds.type });
  
  const apiUrlInput = page.getByPlaceholder('https://updown.io/api');
  await expect(apiUrlInput).toBeVisible();
  await expect(apiUrlInput).toHaveValue('https://updown.io/api');
});

test('should show API Key field', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<UpdownDataSourceOptions, UpdownSecureJsonData>({ fileName: 'datasources.yml' });
  await createDataSourceConfigPage({ type: ds.type });
  
  const apiKeyInput = page.getByPlaceholder('Enter your UpDown.io API key');
  await expect(apiKeyInput).toBeVisible();
});
