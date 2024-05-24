import { Responses } from '../common/API_Responses';
const S3 = require('../common/S3');

const bucket = process.env.bucketName;

exports.handler = async (event) => {
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
        return Responses._200({ newData });
    } catch (error) {
        console.log('error in S3 write', error);
        return Responses._400({ message: 'Failed to write data by filename' });
    }
};

