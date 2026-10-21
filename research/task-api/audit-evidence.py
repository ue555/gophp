#!/usr/bin/env python3
"""Audit raw HTTP artifacts, strict JSON types and complete paired comparisons."""
import argparse
import json
from pathlib import Path
import re

from verify import VERSIONS, check, normalize


def canonical(value):
    # Unlike Python ==, this distinguishes true from 1 and false from 0.
    return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(',', ':'))


def audit(directory, versions=None):
    versions = versions or list(VERSIONS)
    for version in versions:
        paired = []
        for runtime in ('official', 'go'):
            records = [json.loads(line) for line in (directory/f'{runtime}-{version}-http.jsonl').read_text().splitlines()]
            check(len(records) == 191, f'incomplete scenario: {runtime}-{version}')
            comparable = []
            for record in records:
                status = record['status']
                check(status in (200,201,204,404,422), f'unexpected HTTP status: {record}')
                if status == 204:
                    check(record['raw'] == '' and record['content_type'] is None, '204 must have no body or Content-Type')
                    value = None
                else:
                    check(record['content_type'] == 'application/json', 'unexpected Content-Type')
                    value = json.loads(record['raw'])
                if status in (200,201):
                    items = value if isinstance(value,list) else [value]
                    ids = []
                    for item in items:
                        check(set(item) == {'id','title','completed','created_at','updated_at'}, 'Task fields differ')
                        check(type(item['id']) is int and item['id'] > 0, 'id must be a positive JSON integer')
                        check(type(item['title']) is str and 1 <= len(item['title']) <= 120, 'invalid title')
                        check(type(item['completed']) is bool, 'completed must be a JSON boolean')
                        for key in ('created_at','updated_at'):
                            check(re.fullmatch(r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.000000Z',item[key]), 'timestamp format')
                        ids.append(item['id'])
                    check(ids == sorted(set(ids)), 'list ids not unique/ascending')
                record.pop('raw')
                record['json'] = normalize(value)
                comparable.append(record)
            check(canonical(comparable) == canonical(json.loads((directory/f'{runtime}-{version}-normalized.json').read_text())), 'normalized artifact inconsistent with raw response')
            check(sum(r['path'].startswith('/api/tasks?round=') for r in records) == 100, 'missing sequential requests')
            check(sum(r['path'].startswith('/api/tasks?parallel=') for r in records) == 8, 'missing parallel requests')
            paired.append(comparable)
            print(f'PASS raw artifact audit {runtime}-{version}: 191 responses, exact JSON types')
        check(canonical(paired[0]) == canonical(paired[1]), f'strict paired comparison failed: {version}')
        print(f'PASS strict HTTP comparison official {VERSIONS[version]} / Go {version}')
    print(f'PASS all {382 * len(versions)} real HTTP response artifacts; runtime SQLite assertions are recorded separately in run.log')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('directory', type=Path)
    parser.add_argument('--version', choices=VERSIONS, action='append')
    args = parser.parse_args()
    audit(args.directory, args.version)
