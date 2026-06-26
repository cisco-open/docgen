#!/usr/bin/env node

// Copyright 2026 Cisco Systems, Inc. and their affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/**
 * OOXML Validator CLI
 * 
 * Validates a .docx file using @xarsh/ooxml-validator
 * Exits with code 0 if valid, non-zero if validation errors are found
 */

import { validateFile } from '@xarsh/ooxml-validator';
import fs from 'fs';
import path from 'path';

async function main() {
  // Parse command line arguments
  const args = process.argv.slice(2);
  
  if (args.length === 0) {
    console.error('Usage: node validate.js <path-to-docx-file>');
    process.exit(1);
  }
  
  const filePath = args[0];
  
  // Check if file exists
  if (!fs.existsSync(filePath)) {
    console.error(`Error: File not found: ${filePath}`);
    process.exit(2);
  }
  
  // Check if it's a file (not a directory)
  const stats = fs.statSync(filePath);
  if (!stats.isFile()) {
    console.error(`Error: Not a file: ${filePath}`);
    process.exit(2);
  }
  
  console.log(`Validating OOXML document: ${path.basename(filePath)}`);
  
  try {
    // Validate the OOXML document
    const result = await validateFile(filePath);
    
    if (result.ok) {
      console.log('✓ Document is valid OOXML');
      process.exit(0);
    } else {
      console.error('✗ Document validation failed');
      console.error('\nValidation errors:');
      
      if (result.errors && result.errors.length > 0) {
        result.errors.forEach((error, index) => {
          const errorMsg = typeof error === 'string' ? error : (error.message || JSON.stringify(error));
          console.error(`  ${index + 1}. ${errorMsg}`);
        });
      } else {
        console.error('  Unknown validation error');
      }
      
      process.exit(1);
    }
  } catch (error) {
    // Handle errors from the validator itself (e.g., malformed ZIP, corrupted file)
    console.error(`Error during validation: ${error.message}`);
    if (error.code) {
      console.error(`Error code: ${error.code}`);
    }
    process.exit(3);
  }
}

main();
