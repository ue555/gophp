<?php
if (true) {
    echo 'then';
} else {
    echo 'bad';
}
if (false) {
    echo 'bad';
} else {
    echo '|else';
}
if (1 === '1') {
    echo '|bad';
} else {
    echo '|strict';
}
echo "\n";
