// Measure the page load time
var loadTimeStart = new Date().getTime();

// Function to calculate and display the load time
function showLoadTime(label) {
    var loadTimeEnd = new Date().getTime();
    var loadTime = loadTimeEnd - loadTimeStart;
    var loadTimeElement = document.getElementById("load-time");

    // Display the load time in the specified element
    if (loadTimeElement) {
        loadTimeElement.innerHTML+="\t";
        if (loadTimeElement.innerHTML == "Loading...\t") {
            loadTimeElement.innerHTML = "";
        }
        loadTimeElement.innerHTML += "\t"+ label + loadTime + "ms";
    }
}

document.addEventListener("DOMContentLoaded", function(){
    showLoadTime("\nDCL: "); // DOM content loaded
});

// Attach the showLoadTime function to the window's load event
window.addEventListener("load", function(){
    showLoadTime("L: ");
    var performanceTiming = window.performance.timing;
    var firstPaint = performanceTiming.responseStart - performanceTiming.navigationStart;
    showLoadTime("FP: "); // First paint
    var domInteractive = performanceTiming.domInteractive - performanceTiming.navigationStart;
    showLoadTime("I: "); // DOM interactive
    var domComplete = performanceTiming.domComplete - performanceTiming.navigationStart;
    // showLoadTime("DOM complete: ");
    // var firstByte = performanceTiming.responseStart - performanceTiming.navigationStart;
});