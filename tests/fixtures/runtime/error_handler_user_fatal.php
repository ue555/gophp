<?php
echo "before|";
try {
    trigger_error('fatal', E_USER_ERROR);
} finally {
    echo "unreachable\n";
}
