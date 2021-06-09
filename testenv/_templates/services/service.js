export class {{ .Name }}Service {
  name() {
    return "{{ .Name }}";
  }

  serve() {
    return "ok";
  }
}
