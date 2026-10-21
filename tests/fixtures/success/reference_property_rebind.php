<?php
class RefPair { public $left = 1; public $right = 2; }
$pair = new RefPair;
$old =& $pair->left;
$other =& $pair->left;
$pair->left =& $pair->right;
$pair->right = 9;
echo $old, '|', $other, '|', $pair->left, '|', $pair->right;
