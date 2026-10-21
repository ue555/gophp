<?php
function loadDuplicate() {
    echo 'before|';
    try {
        require __DIR__ . '/../support/include_integration_duplicate.inc';
    } finally {
        echo 'finally|';
        return 9;
    }
}
echo loadDuplicate(), '|after';
