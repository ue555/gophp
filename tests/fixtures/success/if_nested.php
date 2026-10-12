<?php
IF (1 === 1) {
    if (false) {
        echo 'bad';
    } ELSE /* between blocks */ {
        if ('00') { echo 'deep'; } else { echo 'bad'; }
    }
    echo '|outer';
} else {
    echo 'bad';
}
if (false) {
    echo 'bad';
} else {
    if ('0') { echo 'bad'; } else { echo '|else-nested'; }
}
echo "\n";
