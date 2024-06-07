import { Responses } from '../common/API_Responses';
const S3 = require('../common/S3');
const DynamoDB = require("../common/Dynamo");

const bucket = process.env.bucketName;
const table = process.env.tableName;

exports.handler = async (event: any) => {
    console.log('event', event);

    if (!event.pathParameters || !event.pathParameters.fileName) {
        return Responses._400({ message: 'missing the fileName from the path' });
    }

    if (!event.pathParameters || !event.pathParameters.hash) {
        return Responses._400({ message: 'missing the hash from the path' });
    }

    let fileName = event.pathParameters.fileName;
    let hash = event.pathParameters.hash;

    try {
        const buffer = Buffer.from(event.body, 'base64'); // Ensure the body is interpreted correctly
        const newData = await S3.write(buffer, fileName, bucket, hash);
        
        // Save the hash and filename to DynamoDB
        await DynamoDB.write(hash, fileName, table);

        return Responses._200({ newData });
    } catch (error) {
        console.log('error in S3 write', error);
        return Responses._400({ message: 'Failed to write data by filename' });
    }
};