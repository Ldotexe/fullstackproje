const { protocol, hostname, port } = window.location;

const apiURL = port ? `${protocol}//${hostname}:${port}/api` : `${protocol}//${hostname}/api`;
export { apiURL };