<?php
$alias =& $missing;
echo gettype($missing), '|';
$alias = 5;
echo $missing, '|';
$self =& $self;
echo gettype($self), '|';
$self = 8;
echo $self;
