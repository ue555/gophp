<?php
class TypedDetach {
    public int $first = 1;
    public int $second = 2;
}
$box = new TypedDetach;
$alias =& $box->first;
$box->second =& $alias;
unset($box->first);
$alias = 4;
echo $box->second, '|';
unset($box->second);
$alias = 'free';
echo $alias, '|';
$box->second = 7;
echo $box->second;
