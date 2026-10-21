<?php
try {
    echo 'before|';
    exit('stop');
} finally {
    echo 'wrong';
}
