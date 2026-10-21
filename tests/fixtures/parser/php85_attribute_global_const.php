<?php
set_error_handler(function ($level, $message) {
    echo $level, '|';
});
#[Deprecated(message: 'old')]
const OLDER = 4;
echo OLDER;
