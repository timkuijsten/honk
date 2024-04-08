
function hotkey(e) {
	if (e.ctrlKey || e.altKey)
		return
	if (e.code == "Escape") {
		var menu = document.getElementById("topmenu")
		menu.open = false
		return
	}
	if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement)
		return

	switch (e.code) {
	case "KeyR":
		refreshhonks(document.getElementById("honkrefresher"));
		break;
	case "KeyS":
		oldestnewest(document.getElementById("newerscroller"));
		break;
	case "KeyJ":
		scrollnexthonk();
		break;
	case "KeyK":
		scrollprevioushonk();
		break;
	case "KeyM":
		var menu = document.getElementById("topmenu")
		if (!menu.open) {
			menu.open = true
			menu.querySelector("a").focus()
		} else {
			menu.open = false
		}
		break
	case "Slash":
		document.getElementById("topmenu").open = true
		document.getElementById("searchbox").focus()
		e.preventDefault()
		break
	}
}

(function() {
	document.addEventListener("keydown", hotkey)
	var totop = document.querySelector(".nophone")
	if (totop) {
		totop.onclick = function() {
			window.scrollTo(0,0)
		}
	}
	var els = document.getElementsByClassName("donklink")
	while (els.length) {
		let el = els[0]
		el.onclick = function() {
			el.classList.remove("donk")
			el.onclick = null
			return false
		}
		el.classList.remove("donklink")
	}

})()
