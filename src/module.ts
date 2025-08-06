import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { UpdownQuery, UpdownDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, UpdownQuery, UpdownDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
