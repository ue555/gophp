<?php
$a = '0'; echo (++$a === 1) + 0;
$a = ' 2 '; echo (--$a === 1) + 0;
$a = ' '; echo (++$a === ' ') + 0;
$a = '12tail'; echo (++$a === '12taim') + 0, "\n";
