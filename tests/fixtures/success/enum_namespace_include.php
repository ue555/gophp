<?php
require __DIR__ . '/../support/enum_namespaced.inc';
echo \Catalog\Stage::Ready->name, '|';
echo \Catalog\Stage::from('done')->value;
