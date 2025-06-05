// Busca a URL do asset pelo nome na release
func (c *Client) GetReleaseAssetURL(pkg, version, assetName string) (string, error) {
    // Busca release no GitLab pela tag version (ex: v1.2.3), lista assets e procura assetName
    // Retorna a URL do asset (download direto)
    // ...
}