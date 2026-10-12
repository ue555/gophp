<?php
echo 'before|';
if (true) {
    if (false) { echo 'bad'; } else {
        echo 1 / 0;
        echo 'bad';
    }
}
echo 'after';
