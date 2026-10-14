<?php
for ($i = 0; $i < 3; $i++) {
    switch ($i):
        case 0:
            echo "a";
        case 1:
            echo "b";
            break;
        default:
            echo "c";
    endswitch;
}
echo "\n";
