<?php
if (false) {
    if ($missing / 0) { echo 'bad'; } else { echo 'bad'; }
}
if (true) {
    if (false) {
        if (1 / 0) { echo 'bad'; }
    } else {
        echo 'selected';
    }
} else {
    if ($alsoMissing) { echo 'bad'; }
}
echo "\n";
