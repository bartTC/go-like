/*
    <div class="like-button" data-id="foo" data-domain="example.com">
        <span class="like-button-text">I like this!</span>
        <span class="like-button-counter">0</span>
    </div>
*/


// Retrieve the like count for every object on the page
var likeButtonApiURL = 'http://like.mahner.org/',
    buttonSelector = '.like-button',
    buttons = document.querySelectorAll(buttonSelector);

function getButtonData(el, id, push) {
    var method = push ? 'POST' : 'GET',
        url = likeButtonApiURL + id,
        that = this;

    request = new XMLHttpRequest;
    request.open(method, url, false);

    request.onload = function() {
        if (this.status >= 200 && this.status < 400){
            var data = JSON.parse(this.response);
            el.querySelectorAll('.like-button-counter')[0].textContent = data.counter;
        } else {
            console.log('request ok, bad response', this.response);
        }
    };

    request.onerror = function(e) {
        console.log('request error', e);
    };

    request.send();
}

document.addEventListener('DOMContentLoaded', function(){
    Array.prototype.forEach.call(buttons, function(el, i){
        var el = buttons[i],
            id = el.getAttribute('data-id');

        // Update button data
        getButtonData(el, id);

        // Allow one click per element
        el.addEventListener('click', function(){
            getButtonData(el, id, true);
            el.removeEventListener('click', arguments.callee, false);
        });
    });
});
