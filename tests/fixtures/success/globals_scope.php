<?php
$counter = 1;
function updateGlobal(): void {
    $GLOBALS['counter'] = 2;
    $copy = $GLOBALS;
    $copy['counter'] = 9;
    $GLOBALS['memo']['id'] = true;
}
updateGlobal();
echo $counter, '|', $GLOBALS['counter'], '|';
echo isset($GLOBALS['memo']['id']) ? 'present' : 'missing';
