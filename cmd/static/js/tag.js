


export function positionSuggestionBox(inputElement, suggestionBox) {
    let rect = inputElement[0].getBoundingClientRect();
    let caretPosition = getCaretCoordinates(inputElement);

    suggestionBox.css({
        left: rect.left + caretPosition.left + 'px',
        top: rect.top + caretPosition.top + $(inputElement).scrollTop() + 25 + 'px',
        display: 'block'
    });
}

export function insertAtCaret(elementId, text) {
    const textarea = document.getElementById(elementId);
    let start = textarea.selectionStart;
    let end = textarea.selectionEnd;
    let value = textarea.value;

    // Find the start of the mention to be replaced
    let mentionStart = value.lastIndexOf('@', start - 1);
    if (mentionStart === -1 || mentionStart < value.lastIndexOf(' ', start - 1)) {
        mentionStart = start;
    }
    
    // Replace the text from the start of the mention to the current caret position
    let newValue = value.substring(0, mentionStart) + text + value.substring(end);

    // Update the textarea value and caret position
    textarea.value = newValue;
    textarea.selectionStart = textarea.selectionEnd = mentionStart + text.length;

    textarea.focus();
}

function getCaretCoordinates(element) {
    const rect = element[0].getBoundingClientRect();
    let caretPos = element.selectionStart;
    let fontSize = parseInt(window.getComputedStyle(element[0]).fontSize, 10);
    let top = rect.top + Math.floor(caretPos / element[0].cols) * fontSize;
    let left = rect.left + (caretPos % element[0].cols) * fontSize * 0.6;

    return { top: top, left: left };
}