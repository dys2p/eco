package captcha

type TemplateData struct {
	Answer string // old answer, maybe incorrect
	Err    bool
	ID     string
}

const TemplateString = `
	{{define "captcha"}}
		{{if .}}
			<script>
				function reloadCaptcha() {
					var e = document.getElementById("captcha-image");
					var q = "reload=" + (new Date()).getTime();
					var src  = e.src;
					var p = src.indexOf('?');
					if (p >= 0) {
						src = src.substr(0, p);
					}
					e.src = src + "?" + q
				}
			</script>
			<p class="text-center">
				<img id="captcha-image" src="/captcha/{{.ID}}.png" alt="Captcha image">
			</p>
			<div class="mb-3">
				<label for="captcha-answer" class="form-label">Please solve the captcha:</label>
				<input class="form-control {{if .Err}}is-invalid{{end}}" id="captcha-answer" name="captcha-answer" value="{{.Answer}}" type="number" required>
				<div class="invalid-feedback">Please type the digits correctly.</div>
				<div class="form-text"><a class="text-muted" href="#" onclick="reloadCaptcha(); return false;">Load other captcha image (requires JavaScript)</a></div>
			</div>
			<input type="hidden" name="captcha-id" value="{{.ID}}">
		{{end}}
	{{end}}`
