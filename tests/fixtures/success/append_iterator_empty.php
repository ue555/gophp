<?php
$combined = new AppendIterator;
foreach ($combined as $value) echo 'bad';
echo is_null($combined->getInnerIterator()) ? 'none|' : 'bad|';
$combined->append(new ArrayIterator([4]));
foreach ($combined as $value) echo $value;
echo '|';
foreach ($combined as $value) echo $value;
