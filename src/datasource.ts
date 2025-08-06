import { DataSourceInstanceSettings, CoreApp, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { UpdownQuery, UpdownDataSourceOptions, UpdownCheck, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<UpdownQuery, UpdownDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<UpdownDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<UpdownQuery> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: UpdownQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      checkToken: getTemplateSrv().replace(query.checkToken, scopedVars),
    };
  }

  filterQuery(query: UpdownQuery): boolean {
    // if no query has been provided, prevent the query from being executed
    return !!query.queryType;
  }

  /**
   * Récupère la liste des checks depuis l'API UpDown.io
   */
  async getChecks(): Promise<UpdownCheck[]> {
    return this.getResource('checks');
  }

  /**
   * Teste la connexion à l'API UpDown.io
   */
  async testDatasource() {
    return this.getResource('health');
  }
}
