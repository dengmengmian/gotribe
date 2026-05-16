export interface SystemConfig {
  system_config_id: string
  title: string
  content: string
  logo: string
  icon: string
  footer: string
}

export interface ConfigResponse {
  systemConfig: SystemConfig
}
