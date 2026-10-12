<?php
echo 'before|';
echo (1 / 0) || $must_not_evaluate;
echo 'after';
