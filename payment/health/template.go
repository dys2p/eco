package health

type Item struct {
	Name   string
	Synced bool
}

// workaround because golang can't escape backticks within backticks
const backtick = "`"

const TemplateString = `
	{{define "health"}}
		<script>
			function updateHealth() {
				var xhr = new XMLHttpRequest();
				xhr.onreadystatechange = function() {
					if(xhr.readyState == 4) {
						let content = "";
						if(xhr.status == 200) {
							let data = JSON.parse(xhr.responseText);
							for(item of data) {
								content = content + ` + backtick + `<span class="badge bg-${item.Synced ? 'success' : 'warning'} my-2">${item.Name}: ${item.Synced ? 'synced' : 'out of sync'}</span> ` + backtick + `;
							}
						} else {
							content = '<span class="badge bg-warning">could not connect</span>';
						}
						document.getElementById("health-widget").innerHTML = content;
					}
				}
				xhr.open("GET", "/payment-health", true); // true for asynchronous
				xhr.send(null);

				setTimeout(function(){
					updateHealth();
				}, 10*1000); // 10 seconds
			}
			document.write('<span id="health-widget"></span>');
			updateHealth();
		</script>
		<noscript>
			<span class="badge bg-secondary">requires JavaScript</span>
		</noscript>
	{{end}}`
