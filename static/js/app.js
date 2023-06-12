var tribute = new Tribute({
    values: function (text, cb) {
        remoteSearch(text, users => cb(users));
    },
    allowSpaces: true,
    menuShowMinLength: 3,
});

function remoteSearch(text, cb) {
    var URL = "/students/search/at/";
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                var data = JSON.parse(xhr.responseText);
                cb(data);
            } else if (xhr.status === 403) {
                cb([]);
            }
        }
    };
    xhr.open("GET", URL + "?q=" + text, true);
    xhr.send();
}

// After htmx swaps the content, re-attach tribute to the new textarea
htmx.on('htmx:afterSwap', function (evt) {
    tribute.attach(document.querySelectorAll('textarea'));
    var textareas = document.querySelectorAll('textarea');
    for (var i = 0; i < textareas.length; i++) {
        textarea = textareas[i];
        textarea.style.height = textarea.scrollHeight + "px";
        textarea.oninput = function () {
            if (textarea.style.height != textarea.scrollHeight + "px") {
                textarea.style.height = textarea.scrollHeight + "px";
            }
        };
    }
    
});

document.addEventListener('keydown', function(e) {
    if (document.activeElement.tagName == 'INPUT' && e.key == 'Escape') {
        var form = document.activeElement.closest('form');
        if (form) {
            var submit = form.querySelector('input[type="button"]');
            if (submit) {
                submit.click();
            }
        }
    }
    if (e.ctrlKey && e.key == '.') {
        e.preventDefault();
        var sel = window.getSelection();
        
    }
    if (e.ctrlKey && e.key == 'f') {
        search = document.getElementById('search');
        if (search) {
            search.focus();
            e.preventDefault();
        }
    }
    if (document.activeElement.tagName != 'TEXTAREA') {
        return;
    }
    if (e.key == 'Tab') {
        e.preventDefault();
        var textarea = document.activeElement;
        var start = textarea.selectionStart;
        var end = textarea.selectionEnd;
        var text = textarea.value;
        var selectedText = text.substring(start, end);
        var newText = text.substring(0, start) + "    " + text.substring(end);
        textarea.value = newText;
        textarea.selectionStart = start + 4;
        textarea.selectionEnd = end + 4;
    }
    if (e.ctrlKey && (e.key == 'b' || e.key == 'i' || e.key == 'u' || e.key == 'l' || e.key == '>' || e.key == 'Enter') || e.key == 'Escape') {
        var keyCode = e.key;
        var focused = document.activeElement;
        var id = focused.id;
        e.preventDefault();
        if (keyCode == 'b') {
            insertFormating(focused, "**", "bold", "**");
        } else if (keyCode == 'i') {
            insertFormating(focused, "_", "italic", "_");
        } else if (keyCode == 'u') {
            insertFormating(focused, "++", "underline", "++");
        } else if (keyCode == 'l') {
            insertFormating(focused, "[", "link title", "](http://www.example.com)");
        } else if (keyCode == '>') {
            insertFormating(focused, "> ", "quote", "");
        } else if (keyCode == 'Enter') {
            // Trigger the submit button in the form
            var form = focused.closest('form');
            if (form) {
                var submit = form.querySelector('input[type="submit"]');
                if (submit) {
                    submit.click();
                }
            }
        } else if (keyCode == 'Escape') {
            // Trigger the cancel button in the form
            var form = focused.closest('form');
            if (form) {
                var submit = form.querySelector('input[type="button"]');
                if (submit) {
                    submit.click();
                }
            }
        }
    }
});

function insertFormating(txtarea, text, defaultTxt = "", text2 = "") {
    var selectStart = txtarea.selectionStart
    var selectEnd = txtarea.selectionEnd
    var scrollPos = txtarea.scrollTop;
    var caretPos = txtarea.selectionStart;
    var mode = 0; //Adding markdown with selected text
    var front = (txtarea.value).substring(0, caretPos);
    var back = (txtarea.value).substring(selectEnd, txtarea.value.length);
    var middle = (txtarea.value).substring(caretPos, selectEnd);
    
    var textLen = text.length;
    var text2Len = text2.length;
    
    if (selectStart === selectEnd) {
        middle = defaultTxt;
        mode = 1; //Adding markdown with default text
    } else {
        if (front.substr(-textLen) == text && back.substr(0, text2Len) == text2) {
            front = front.substring(0, front.length - textLen);
            back = back.substring(text2Len, back.length);
            text = "";
            text2 = "";
            mode = 2; //Removing markdown with selected text eg. **<selected>bold<selected>**
        } else if (middle.substr(0, textLen) == text && middle.substr(-text2Len) == text2) {
            middle = middle.substring(textLen, middle.length - text2Len);
            text = "";
            text2 = "";
            mode = 3; //Removing markdown with selected text eg. <selected>**bold**<selected>
        }
    }
    txtarea.value = front + text + middle + text2 + back;
    if (selectStart !== selectEnd) {
        if (mode === 0) {
            txtarea.selectionStart = selectStart + textLen;
            txtarea.selectionEnd = selectEnd + textLen;
        } else if (mode === 2) {
            txtarea.selectionStart = selectStart - textLen;
            txtarea.selectionEnd = selectEnd - textLen;
        } else if (mode === 3) {
            txtarea.selectionStart = selectStart;
            txtarea.selectionEnd = selectEnd - textLen - text2Len;
        }
    } else {
        txtarea.selectionStart = selectStart + textLen;
        txtarea.selectionEnd = txtarea.selectionStart + middle.length;
    }
    txtarea.focus();
    txtarea.scrollTop = scrollPos;
}
/*
* Modal
*
* Pico.css - https://picocss.com
* Copyright 2019-2023 - Licensed under MIT
*/

// Config
const isOpenClass = 'modal-is-open';
const openingClass = 'modal-is-opening';
const closingClass = 'modal-is-closing';
const animationDuration = 400; // ms
let visibleModal = null;


// Toggle modal
const toggleModal = event => {
    event.preventDefault();
    if (event.currentTarget.id == 'confirm-delete') {
        document.getElementById(event.currentTarget.getAttribute('data-delete')).click();
    }
    const modal = document.getElementById(event.currentTarget.getAttribute('data-target'));
    document.getElementById('confirm-delete').setAttribute('data-delete', event.currentTarget.getAttribute('data-delete'));
    (typeof(modal) != 'undefined' && modal != null)
    && isModalOpen(modal) ? closeModal(modal) : openModal(modal)
}

// Is modal open
const isModalOpen = modal => {
    return modal.hasAttribute('open') && modal.getAttribute('open') != 'false' ? true : false;
}

// Open modal
const openModal = modal => {
    if (isScrollbarVisible()) {
        document.documentElement.style.setProperty('--scrollbar-width', `${getScrollbarWidth()}px`);
    }
    document.documentElement.classList.add(isOpenClass, openingClass);
    setTimeout(() => {
        visibleModal = modal;
        document.documentElement.classList.remove(openingClass);
    }, animationDuration);
    modal.setAttribute('open', true);
}

// Close modal
const closeModal = modal => {
    visibleModal = null;
    document.documentElement.classList.add(closingClass);
    setTimeout(() => {
        document.documentElement.classList.remove(closingClass, isOpenClass);
        document.documentElement.style.removeProperty('--scrollbar-width');
        modal.removeAttribute('open');
    }, animationDuration);
}

// Close with a click outside
document.addEventListener('click', event => {
    if (visibleModal != null) {
        const modalContent = visibleModal.querySelector('article');
        const isClickInside = modalContent.contains(event.target);
        !isClickInside && closeModal(visibleModal);
    }
});

// Close with Esc key
document.addEventListener('keydown', event => {
    if (event.key === 'Escape' && visibleModal != null) {
        closeModal(visibleModal);
    }
});

// Get scrollbar width
const getScrollbarWidth = () => {
    
    // Creating invisible container
    const outer = document.createElement('div');
    outer.style.visibility = 'hidden';
    outer.style.overflow = 'scroll'; // forcing scrollbar to appear
    outer.style.msOverflowStyle = 'scrollbar'; // needed for WinJS apps
    document.body.appendChild(outer);
    
    // Creating inner element and placing it in the container
    const inner = document.createElement('div');
    outer.appendChild(inner);
    
    // Calculating difference between container's full width and the child width
    const scrollbarWidth = (outer.offsetWidth - inner.offsetWidth);
    
    // Removing temporary elements from the DOM
    outer.parentNode.removeChild(outer);
    
    return scrollbarWidth;
}

// Is scrollbar visible
const isScrollbarVisible = () => {
    return document.body.scrollHeight > screen.height;
}

var dropZone = document.getElementById('dropZone');

function showDropZone() {
    dropZone.style.visibility = "visible";
}
function hideDropZone() {
    dropZone.style.visibility = "hidden";
}

function allowDrag(e) {
    if (true) {  // Test that the item being dragged is a valid one
        e.dataTransfer.dropEffect = 'copy';
    }
}

function handleDrop(e) {
    hideDropZone();
}

// 1
window.addEventListener('dragenter', function(e) {
    dropZone = document.getElementById('dropZone');
    // 2
    dropZone.addEventListener('dragenter', allowDrag);
    
    // 3
    dropZone.addEventListener('dragleave', function(e) {
        hideDropZone();
    });
    
    // 4
    dropZone.addEventListener('drop', handleDrop);
    if (dropZone == null) return;
    showDropZone();
});


// Copy the inner text of a node with class "copy" to the clipboard when clicked
document.addEventListener('click', async function(e) {
    if (e.target.classList.contains('copy')) {
        if (e.target.hasAttribute('data-value')) {
            var text = e.target.getAttribute('data-value');
        } else {
            var text = e.target.innerText;
        }
        try {
            await navigator.clipboard.writeText(text);
            e.target.setAttribute('data-tooltip', 'Copied!');
            setTimeout(function() {
                e.target.setAttribute('data-tooltip', 'Click to copy');
            }, 1000);
        } catch (err) {
            console.error('Failed to copy: ', err);
        }
    }
});
