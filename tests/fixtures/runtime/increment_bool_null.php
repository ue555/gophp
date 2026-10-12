<?php
$a = true; echo ($a++ === true) + 0, ($a === true) + 0;
$a = false; echo (--$a === false) + 0;
$a = null; echo ($a-- === null) + 0, ($a === null) + 0, "\n";
