'use strict'

module.exports = test

const assert = require('assert')

function test({workspaces, tags, variables}) {
  const testCases = [
    {
      test: 'workspaces output',
      message: 'workspaces array does not equal expected output',
      actual: workspaces,
      expected: ['staging', 'production'],
    },
    {
      test: 'workspace_tags output',
      message: 'workspace tag map does not equal the expected output',
      actual: tags,
      expected: {
        staging: ['environment:staging'],
        production: ['environment:production'],
      },      
    },
    {
      test: 'workspace_variables output',
      message: 'workspace variables map does not equal the expected output',
      actual: variables,
      expected:{
        staging: [{
          key: 'environment',
          value: 'staging',
          category: 'terraform',
        }],
        production: [{
          key: 'environment',
          value: 'production',
          category: 'terraform',      
        }]
      }  
    }
  ]

  testCases.forEach(({test, actual, expected, message}) => {
    console.log(`${test}...`)
    console.group()
    assert.deepEqual(actual, expected, message)
    console.log('OK')
    console.groupEnd()
  })
}
