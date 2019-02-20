

(function() {
	var insert = document.querySelector(
			"#AlgorithmSpecificControls>td>input[type=Text]"
		),
		sleep = function(ms) {
			return new Promise(resolve => setTimeout(resolve, ms));
		},
		array = [
			57,11,7,0,87,75,21,47,51,90,36,88,69,70,97,66,20,74,30,8,35,12,34,
			52,25,81,92,79,26,49,2,39,58,23,85,53,99,86,31,48,64,54,71,55,60,72,
			78,4,62,93,44,100,77,89,13,29,83,80,59,32,68,18,43,3,38,5,95,96,10,
			1,37,50,63,15,82,76,22,67,45,28,16,17,56,84,24,40,91,33,19,73,46,27,
			9,14,6,98,41,42,65,94,61
		];

	(async function(){
		for (var j = 0; j < array.length; j++) {
			var i = array[j];
			console.log("insert", i)
			insert.value = i;
			insert.dispatchEvent(
				new KeyboardEvent("keydown", {
					bubbles: true,
					cancelable: true,
					which: 13
				})
			);
			await sleep(200);
			for (;insert.disabled;) {
				await sleep(200);
			}
		}
	})()
})()
