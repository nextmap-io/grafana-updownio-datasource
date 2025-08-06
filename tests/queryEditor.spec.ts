import { test, expect } from '@grafana/plugin-e2e';

test('smoke: should render query editor', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  
  // Check that the query editor renders by looking for the select dropdown
  const queryEditorRow = panelEditPage.getQueryEditorRow('A');
  await expect(queryEditorRow.locator('[class*="react-select"]')).toBeVisible();
});

test('should show default query type selection', async ({
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  
  // Should show the Select data type dropdown placeholder
  const queryEditorRow = panelEditPage.getQueryEditorRow('A');
  await expect(queryEditorRow.getByText('Select data type')).toBeVisible();
});

test('should render query editor interface', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  await panelEditPage.setVisualization('Table');
  
  // Just verify the interface loads by checking for react-select component
  const queryEditorRow = panelEditPage.getQueryEditorRow('A');
  await expect(queryEditorRow.locator('[class*="react-select"]')).toBeVisible();
});
