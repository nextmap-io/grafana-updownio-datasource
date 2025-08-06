import { test, expect } from '@grafana/plugin-e2e';

test('smoke: should render query editor', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  await expect(panelEditPage.getQueryEditorRow('A').getByText('Data type')).toBeVisible();
});

test('should show service selector when metrics is selected', async ({
  panelEditPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  
  // Select metrics query type
  await panelEditPage.getQueryEditorRow('A').getByText('Select data type').click();
  await page.getByText('Service Metrics').click();
  
  // Service selector should appear
  await expect(panelEditPage.getQueryEditorRow('A').getByText('Service')).toBeVisible();
});

test('should render service list query by default', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  await panelEditPage.setVisualization('Table');
  await expect(panelEditPage.refreshPanel()).toBeOK();
});
